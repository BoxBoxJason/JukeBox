import pandas as pd
from sklearn.model_selection import train_test_split
from transformers import DistilBertForSequenceClassification, Trainer, TrainingArguments, DistilBertTokenizerFast
from torch.utils.data import Dataset, DataLoader

MODEL_PATH = "moderation_model"
PRETRAINED_DISTILBERT_MODEL = "distilbert-base-uncased"

def openDataset(path):
    """
    openDataset function reads the dataset from the given path and returns a pandas DataFrame.
    : path: str: The path to the dataset file.
    : return: pd.DataFrame: The dataset as a pandas DataFrame.
    """
    return pd.read_csv(path)


def preProcessData(data):
    """
    preProcessData function takes a pandas DataFrame and returns the inputs, labels, and tokenizer.
    : data: pd.DataFrame: The dataset as a pandas DataFrame.
    : return: Tuple[Dict[str, torch.Tensor], np.ndarray, DistilBertTokenizerFast]: A tuple containing the inputs, labels, and tokenizer.
    """
    tokenizer = DistilBertTokenizerFast.from_pretrained(PRETRAINED_DISTILBERT_MODEL)
    inputs = tokenizer(list(data['text']), truncation=True, padding=True, max_length=128)
    labels = data['label'].values

    return inputs, labels, tokenizer


# Split the data into training and validation sets
class ModerationDataset(Dataset):
    def __init__(self, encodings, labels):
        self.encodings = encodings
        self.labels = labels

    def __getitem__(self, idx):
        item = {key: torch.tensor(val[idx]) for key, val in self.encodings.items()}
        item['labels'] = torch.tensor(self.labels[idx])
        return item

    def __len__(self):
        return len(self.labels)


def splitDataset(inputs, labels):
    """
    splitDataset function takes the inputs and labels and returns the training and validation datasets.
    : inputs: Dict[str, torch.Tensor]: The inputs for the model.
    : labels: np.ndarray: The labels for the inputs.
    : return: Tuple[ModerationDataset, ModerationDataset]: A tuple containing the training and validation datasets.
    """
    train_encodings, val_encodings, train_labels, val_labels = train_test_split(inputs, labels, test_size=0.2)
    train_dataset = ModerationDataset(train_encodings, train_labels)
    validation_dataset = ModerationDataset(val_encodings, val_labels)
    return train_dataset, validation_dataset


def trainModerationModel(train_dataset, validation_dataset, tokenizer, save_path=None):
    """
    trainModerationModel function takes the training and validation datasets, tokenizer, and an optional save path.
    It trains the moderation model and saves the model and tokenizer to the save path if provided.
    : train_dataset: ModerationDataset: The training dataset.
    : validation_dataset: ModerationDataset: The validation dataset.
    : tokenizer: DistilBertTokenizerFast: The tokenizer for the model.
    : save_path: str: The path to save the model and tokenizer.
    : return: None
    """
    moderation_model = DistilBertForSequenceClassification.from_pretrained(PRETRAINED_DISTILBERT_MODEL, num_labels=2)
    # Training setup
    training_args = TrainingArguments(
        output_dir="./results",
        evaluation_strategy="epoch",
        per_device_train_batch_size=16,
        per_device_eval_batch_size=16,
        num_train_epochs=3,
        weight_decay=0.01,
    )

    trainer = Trainer(
        model=moderation_model,
        args=training_args,
        train_dataset=train_dataset,
        eval_dataset=validation_dataset,
    )

    trainer.train()
    if save_path:
        trainer.save_model(save_path)
        tokenizer.save_pretrained(save_path)
