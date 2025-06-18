Feature: PR Strategy
  As a release engineer
  I want to apply different strategies for committing rendered manifests
  So that I can either push directly to the default branch or open a pull request

  Background:
    Given a repository with default branch "main"
    And no branch named matching "appsync/*"
    And the following rendered files:
      | path                              | content          |
      | .applications/app-1/application.yaml | apiVersion: v1   |
      | .applications/app-1/persistence.yaml | apiVersion: v1   |
      | .applications/app-1/edge.yaml        | apiVersion: v1   |

  Scenario: Direct commit strategy
    Given I select the direct commit strategy
    When I apply the strategy to the rendered files
    Then the files are written to branch "main"
    And no pull request is created

  Scenario: Feature–branch PR strategy
    Given I select the feature-branch PR strategy
    When I apply the strategy to the rendered files
    Then a new branch matching "appsync/\\d+" is created from "main"
    And the files are written to that new branch
    And a pull request titled "appsync sync" is opened from the new branch into "main"
