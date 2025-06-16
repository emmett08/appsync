Feature: Skeleton rollout into tenant repos
  In order to bootstrap 1AI resources
  As platform engineer
  I want to run the CLI and have
  skeleton files transformed into CRDs
  and a pull request opened against each tenant repo

  Background:
    Given a skeleton directory "testdata/skeleton" containing:
      | path               | content                                  |
      | app-1/persistence.yaml | apiVersion: v1\nkind: Persistence\n... |
      | app-1/edge.yaml        | apiVersion: v1\nkind: Edge\n...        |
    And a fake tenant repo at "testdata/repos/ott-ref-app" with a main branch
    And a config file "testdata/config.yaml" containing:
      | skeleton        | testdata/skeleton          |
      | applicationsDir | .applications              |
      | tenantRepos     | testdata/repos/ott-ref-app |
      | tenantRepos     | testdata/repos/ovp         |

  Scenario: Run rollout writes files and opens PR
    When I execute `1ai-pr rollout --config testdata/config.yaml --dry-run`
    Then the directory "testdata/repos/ott-ref-app/.applications/app-1" should contain:
      | file             |
      | persistence.yaml |
      | edge.yaml        |
    And a pull request should have been opened for "ott-ref-app" on branch "1ai/"
