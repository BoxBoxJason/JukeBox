from transformers import T5ForConditionalGeneration, T5Tokenizer

rephrase_tokenizer = T5Tokenizer.from_pretrained("t5-small")
rephrase_model = T5ForConditionalGeneration.from_pretrained("t5-small")

def rephrase(text):
    """
    rephrase function takes a text input and returns a rephrased version of the input.
    : text: str: The text input to be rephrased.
    : return: str: A rephrased version of the input.
    """
    input_text = f"paraphrase: {text} </s>"
    encoding = rephrase_tokenizer.encode_plus(input_text, return_tensors="pt", max_length=128, padding="max_length", truncation=True)
    input_ids, attention_mask = encoding["input_ids"], encoding["attention_mask"]

    outputs = rephrase_model.generate(
        input_ids=input_ids, attention_mask=attention_mask,
        max_length=128,
        num_beams=4,
        repetition_penalty=2.5,
        length_penalty=1.0,
        early_stopping=True
    )

    return rephrase_tokenizer.decode(outputs[0], skip_special_tokens=True)


def moderateAndRephrase(user_input):
    """
    moderateAndRephrase function takes a user input and returns a rephrased version of the input if it is safe to use.
    If the input is not safe to use, it returns an empty string.
    : user_input: str: The user input to be moderated and rephrased.
    : return: Tuple[bool, str]: A tuple containing a boolean value indicating if the input is safe and a rephrased version of the input.
    """

    inputs = tokenizer(user_input, return_tensors="pt", truncation=True, padding=True, max_length=128)
    outputs = model(**inputs)
    logits = outputs.logits
    prediction = torch.argmax(logits, dim=-1).item()

    if prediction == 1:
        return False, ""

    return True, rephrase(user_input)
