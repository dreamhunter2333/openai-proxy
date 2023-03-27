import openai

openai.api_key = "xxx"
openai.api_base = "http://localhost:8080/v1"

# create a completion
response = openai.Completion.create(
    model="text-davinci-003",
    prompt="Say this is a test",
    temperature=0,
    max_tokens=7
)

# print the completion
print(response)
