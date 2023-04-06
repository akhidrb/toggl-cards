Feature: Draw Cards from Deck

  Scenario: (3.1) Draw Cards from a Deck (Happy Path)
    Given a user creates a partial deck that is not shuffled with the following cards:
      | AC |
      | 5S |
      | JD |
      | KH |
      | 2C |
      | 7S |
      | JC |
    When the user draws 3 card(s) from the deck
    Then the user should get the following cards:
      | code | value | suit   |
      | 2C   | 2     | CLUBS  |
      | 7S   | 7     | SPADES |
      | JC   | JACK  | CLUBS  |
