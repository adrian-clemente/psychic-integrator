psychic-integrator
==================

This project is a release manager written in Go Lang.

It's inteded to work with GIT repositories and Gradle projects.


Project release flow
==================

    1. Create JIRA ticket to be used for release
    2. Set the version to release mode in order to remove SNAPSHOT
    3. Merge DEVELOP branch into MASTER branch from the project selected
    4. Push the commit merge into the repository
    5. Change branch to DEVELOP branch and increment the version of the project and set the version to SNAPSHOT
    6. Commit the project version change
    7. Push to the repository
    8. Create a JIRA release VERSION with the following name "projectName - releaseVersion"
    9. Update all the JIRA issues that were involved in this release with the previous fixVersion created
    8. Send an email with all the JIRA issues ids and the title of each of them.
    
    
