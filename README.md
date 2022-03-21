# everphone-test-task

This is a test-task with Employee and Gift implemented with such functionality

* Basic operations (Full CRUD except Delete method)
* Dataset upload possibility (processed via queue consumer)

I used Clean Architecture approach building this application as
straightforward as I could, but there could be some architectural errors due to lack of time.

I did not have enough time to prepare some unit tests for my code, but, I guess, it is easy to understand how fast
they can be implemented with such architecture

I also noticed duplicate names present in employee dataset and this app will throw an error in this case,
because I found this case too late to have time to add additional check

## Consumer

Consumer is built to handle uploaded datasets in case they are quite huge following Dependency Injection best practices, allowing us to swap the queue realisation.
Moreover, we can customize the handler for messages, but there will always be a default one, as we can see, I created a special message handler for datasets, leaving DefaultHandler as a variant of NullObject.
There are also some nice approaches
in receiver section - we use goroutines and can scale the receivers count if needed via config file provided.

## How to use

- To start, please install: docker, docker-compose and go

- `'make up'` on root folder run it will start all needed containers as well as tidy go modules
  containers with all required dependencies.
- `'make bash'` to open container bash window
- `'make db-init` to prepare migrations tables
- `'make db-migrate` to migrate needed tables
- `localhost:3333` to access go container from localhost
