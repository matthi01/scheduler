To Do:
- rather than sequential ids, start using uuids
- implement authentication with auth0 and jwt
- implement steps for each task
- implement users table
- split out the code in a more logical way
- rather than all the current code duplication, is there a better way to avoid the duplication in tests and the models
- think about how the deleting requests are gonna work - if a category gets deleted, remove all tasks first. If a task gets deleted, remove all the steps first
- version the api - next version will use graphql
- need to add an order property to task
- add steps under tasks
- add users
- add an order
- add an auto-completion for category when all tasks under the category are completed

thoughts:
- should categories have a completed property? If so... we also need a way to complete all tasks underneath that category