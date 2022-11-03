// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

const randomString = () => (Math.random() + 1).toString(36).substring(7)
const randomPassword = () => randomString() + randomString()
const randomEmail = () => randomString() + "@" + randomString() + ".com"

const login = (email, password) => {
  cy.visit("/.ory/ui/login")
  cy.get('[name="identifier"]').type(email)
  cy.get('[name="password"]').type(password)
  cy.get('[name="method"]').click()
  loggedIn(email)
}

const loggedIn = (email) => {
  cy.visit("/.ory/ui/welcome")
  cy.get('[data-testid="logout"] a').should(
    "have.attr",
    "aria-disabled",
    "false",
  )
  cy.visit("/.ory/ui/sessions")
  cy.get("pre").should("contain.text", email)
}

describe("ory proxy", () => {
  const email = randomEmail()
  const password = randomPassword()
  before(() => {
    cy.clearCookies({ domain: null })
  })

  it("navigation works", () => {
    cy.visit("/.ory/ui/registration")
    cy.get('[data-testid="cta-link"]').click()
    cy.location("pathname").should("eq", "/.ory/ui/login")
  })

  it("should be able to execute registration", () => {
    cy.visit("/.ory/ui/registration")
    cy.get('[name="traits.email"]').type(email)
    cy.get('[name="password"]').type(password)
    cy.get('[name="method"]').click()
    cy.visit("/.ory/ui/welcome")
    loggedIn(email)
  })

  it("should be able to execute login", () => {
    login(email, password)
    cy.request("/anything").should((res) => {
      expect(res.body.headers["Authorization"]).to.not.be.empty
      const token = res.body.headers["Authorization"].replace(/bearer /gi, "")

      cy.task(
        "verify",
        res.body.headers["Authorization"].replace(/bearer /gi, ""),
      ).then((decoded) => {
        expect(decoded.session.identity.traits.email).to.equal(email)
      })
    })
  })

  it("should be able to execute logout", () => {
    login(email, password)
    cy.visit("/.ory/ui/welcome")
    cy.get('[data-testid="logout"] a').should(
      "have.attr",
      "aria-disabled",
      "false",
    )

    cy.get('[data-testid="logout"] a').click()
    cy.location("pathname").should("not.contain", "/ui/")
    cy.request("/anything").should((res) => {
      expect(res.body.headers["Authorization"]).to.be.undefined
    })
  })
})
