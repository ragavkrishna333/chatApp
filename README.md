# chatApp
This is a Chat Application used to communicate with members already available in the database using Go broadcasting.

    Private chat is mapped based on the user code.
    In group chat, the messages are mapped to the group code, which is stored in the database.
    Other operations, such as creating a group, adding members to the group, or creating a new member in the environment, are handled with RESTful APIs using Go.
    Chat operations execute on a WebSocket created in Go.
