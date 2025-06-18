Feature: Catalog scanning
  As an operator
  I want to discover all <team>/<app> entries under the catalog root
  So that I can generate descriptors for each application

  Background:
    Given a catalog directory with the following structure:
      | path                              |
      | team1/appA                        |
      | team1/appB                        |
      | team2/appC                        |

  Scenario: scan without filters
    When I scan the catalog root
    Then I should discover 3 descriptors

  Scenario: scan with team filter
    Given I set filter team to "team1"
    When I scan the catalog root
    Then I should discover 2 descriptors
