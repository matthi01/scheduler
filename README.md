To Do:
- rather than sequential ids, start using uuids
- implement authentication with auth0 and jwt
- implement steps for each task
- implement users table
- split out the code in a more logical way
- rather than all the current code duplication, is there a better way to avoid the duplication in tests and the models
- version the api - next version will use graphql
- need to add an order property to task - not sure how this is gonna work... especially with updating... you would have to change the ordering all tasks in the category if you were to move a task from end to start, should this be handled in individual requests or one large request
- add steps under tasks
- add users
- add an auto-completion for category when all tasks under the category are completed

thoughts:
- should categories have a completed property? If so... we also need a way to complete all tasks underneath that category