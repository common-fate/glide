## Getting Started

In the following documentation we will be keeping track of our internal OpenAPI standards for maintaining and creating Common Fate API's

All models, endpoints, responses should follow this standard to keep uniform functionality.

We use code generation libraries in our Go backend and Typescript/React frontend so we have unison types between them when developing. It is important to follow these standards to ensure correct functionality of our API's

- [Frontend codegen library](https://github.com/anymaniax/orval)
- [Backend codegen library](https://github.com/deepmap/oapi-codegen)

## Creating Models

1. For post requests create a `CreateObject` type for the post request body

## Creating Responses

The format of _List_ responses is standardized to the following format:

```
{
  "nextPage": "string",
  "users": [
    {
      "id": "string",
      "oidSub": "string",
      "email": "string",
      "name": "string",
      "picture": "string",
      "isAdmin": true,
      "createdAt": 0,
      "updatedAt": 0
    }
  ]
}
```
