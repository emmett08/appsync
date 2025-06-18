Feature: Catalog scanning
  As an operator
  I want to discover all <team>/<app> entries under the catalogueue root
  So that I can generate descriptors for each application

  Background:
    Given a catalogue directory with the following structure:
      | path                              |
      | team1/appA                        |
      | team1/appB                        |
      | team2/appC                        |

  Scenario: scan without filters
    When I scan the catalogue root
    Then I should discover 3 descriptors

  Scenario: scan with team filter
    Given I set filter team to "team1"
    When I scan the catalogue root
    Then I should discover 2 descriptors

  Scenario: scan with app filter
    Given I set filter app to "appB"
    When I scan the catalogue root
    Then I should discover 1 descriptors
