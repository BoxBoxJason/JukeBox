from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
import uvicorn
from httpx import AsyncClient
from filter import moderate_and_rephrase

app = FastAPI()

MUSIC_GENERATOR_URL = "http://rave:8000"

@app.post("/api/user/message")
async def processInput(request: Request):
    data = await request.json()
    user_input = data.get("message", "")
    flagged, response = moderate_and_rephrase(user_input)

    if flagged:
        # Set the http status to 400 if the message is not safe to use
        return JSONResponse(status_code=400)
    else:
        # Send the sanitized input to the music generator
        await requestToMusicGenerator(response)
        return JSONResponse(status_code=200)

async def requestToMusicGenerator(sanitized_input):
    headers = {
        "Content-Type": "application/json"
    }
    data = {
        "message": sanitized_input
    }
    async with AsyncClient() as client:
        response = await client.post(MUSIC_GENERATOR_URL, headers=headers, json=data)
        return response.json()

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8022)
