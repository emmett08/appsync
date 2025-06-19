Feature: fetch-repos CLI
  In order to maintain a repository list
  As a user
  I want to generate a YAML manifest from the top-level directories in a GitHub repo

  Background:
    Given a running API server that returns the following JSON:
    """
    [
      { "type": "dir",  "team": "teamA", "owner": "ownerA", "repo": "repoA" },
      { "type": "dir",  "team": "teamB", "owner": "ownerB", "repo": "repoB" },
      { "type": "file", "name": "README.md" }
    ]
    """

  Scenario: generate repos.yaml with default regex
    And an output file path "repos.yaml"
    When I run the fetch-repos command for owner "emmett08" repo "appsync"
    Then the exit status should be 0
    And the file "repos.yaml" should contain:
      """
      repos:
        - team: teamA
          owner: ownerA
          repo: repoA
        - team: teamB
          owner: ownerB
          repo: repoB
      """

  Scenario: generate with a custom regex
    Given I use regex `(?P<team>[^_]+)_(?P<repo>.+)`
    And an output file path "custom.yaml"
    When I run the fetch-repos command for owner "emmett08" repo "appsync" regex "(?P<team>[^_]+)_(?P<repo>.+)"
    Then the exit status should be 0
    And the file "custom.yaml" should contain:
      """
      repos:
        - team: teamA
          owner: ownerA
          repo: repoA
        - team: teamB
          owner: ownerB
          repo: repoB
      """

  Scenario: generate repos.yaml from a branch
    Given I use ref "develop"
    And an output file path "branch.yaml"
    When I run the fetch-repos command for owner "emmett08" repo "appsync" ref "develop"
    Then the exit status should be 0
    And the file "branch.yaml" should contain:
      """
      repos:
        - team: teamA
          owner: ownerA
          repo: repoA
        - team: teamB
          owner: ownerB
          repo: repoB
      """

  Scenario: generate repos.yaml from a commit sha
    Given I use ref "deadbeef"
    And an output file path "sha.yaml"
    When I run the fetch-repos command for owner "emmett08" repo "appsync" ref "deadbeef"
    Then the exit status should be 0
    And the file "sha.yaml" should contain:
      """
      repos:
        - team: teamA
          owner: ownerA
          repo: repoA
        - team: teamB
          owner: ownerB
          repo: repoB
      """
