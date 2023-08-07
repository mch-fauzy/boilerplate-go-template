## Contributing

When contributing to this repository, please first discuss the change you wish to make via issue, email, or any other method with the owners of this repository before making a change.

Please note we have a code of conduct, please follow it in all your interactions with the project.

## Pull Request Process
1. Ensure any install or build dependencies are removed before the end of the layer when doing a build.
2. Update the README.md with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations and container parameters.
3. Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent. The versioning scheme we use is SemVer.
4. You may merge the Pull Request in once you have the sign-off of two other developers, or if you do not have permission to do that, you may request the second reviewer to merge it for you.

## Code of Conduct

### Our Pledge
In the interest of fostering an open and welcoming environment, we as contributors and maintainers pledge to making participation in our project and our community a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity and expression, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards
Examples of behavior that contributes to creating a positive environment include:

* Using welcoming and inclusive language
* Being respectful of differing viewpoints and experiences
* Gracefully accepting constructive criticism
* Focusing on what is best for the community
* Showing empathy towards other members

* Examples of unacceptable behavior by participants include:

* The use of sexualized language or imagery and unwelcome sexual attention or advances
* Trolling, insulting/derogatory comments, and personal or political attacks
* Public or private harassment
* Publishing others' private information, such as a physical or electronic address, without explicit permission
* Other conduct which could reasonably be considered inappropriate in a professional setting


## Guidelines
Before we look into how, here are the guidelines. If your Pull Requests fail
to pass these guidelines it will be declined and you will need to re-submit
when you’ve made the changes. This might sound a bit tough, but it is required
for us to maintain quality of the code-base.

## Code Style
All code must meet the Style Guide, SnakeCase, CamelCase or any underscores and readable operators. This makes certain 
that all code is the same format as the existing code and means it will be as readable as possible.

## Documentation
If you change anything that requires a change to documentation then you will need to add it. New packages, methods, 
parameters, changing default values, etc are all things that will require a change to documentation. 
The change-log must also be updated for every change.

## Compatibility
Code should maintain compatibility with the latest version of the code-base.

## Branching
Choose how you want as branching model 

_for example :_ 
"Git Flow" requires all pull requests to be sent to the "develop" branch. This is
where the next planned version will be developed. The "master" branch will always contain the latest stable version and 
is kept clean so a "hotfix" (e.g: an emergency security patch) can be applied to master to create a new version, without 
worrying about other features holding it up. For this reason all commits need to be made to "develop" and any sent to "master" 
will be closed automatically. If you have multiple changes to submit, please place all changes into their own branch on your fork.
One thing at a time: A pull request should only contain one change. That does not mean only one commit, 
but one change - however many commits it took. The reason for this is that if you change X and Y but send a pull 
request for both at the same time, we might really want X but disagree with Y, meaning we cannot merge the request. 
Using the Git-Flow branching model you can create new branches for both of these features and send two requests.

## Keeping your personal branch or fork up-to-date
Always keep your branch or fork up-to-date with the latest stable version of the code-base before doing pull requests.



