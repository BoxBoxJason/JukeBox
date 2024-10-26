# JukeBox

JukeBox is an interactive music creation platform that allows users to create music together by interacting with an AI model in real time.

## AI Models

This project is powered by three distinct AI models:

### Music generation
The music generation AI model is [RAVE](https://forum.ircam.fr/projects/detail/rave-vst/).


This model is an auto encoder, meaning that is takes sound(s) as input and generates a new sound as output. It can generate music in real time based on the user's commands.
It fits our use case perfectly, as it can generate music in real time, and the music's parameters can be controlled by the user's commands.

You can check [their GitHub repository](https://github.com/acids-ircam/RAVE) for more information.

### User requests processing

#### Filtering
The users chat messages are analyzed by an AI model that can filter out messages that are not related to music creation, and can also understand the user's requests and translate them into music commands.
This model's mission is to protect the music generation AI model from being spammed with irrelevant messages. It also recognizes malicious intents and can ban users from the platform.

The filter model is a custom trained version of [DistilBERT](https://huggingface.co/docs/transformers/model_doc/distilbert) trained on a dataset of music-related messages.

#### Reformulation
The reformulation model converts a user request into a music command that can be understood by the music generation AI model.
