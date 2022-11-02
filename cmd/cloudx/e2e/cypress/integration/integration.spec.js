// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

const randomString = () => (Math.random() + 1).toString(36).substring(7)
const randomPassword = () => randomString() + randomString()
const randomEmail = () => randomString() + "@" + randomString() + ".com"

const isTunnel = parseInt(Cypress.env("IS_TUNNEL")) === 1
const prefix = isTunnel ? "" : "/.ory"

const login = (email, password) => {
  cy.visit(prefix + "/ui/login")
  cy.get('[name="identifier"]').type(email)
  cy.get('[name="password"]').type(password)
  cy.get('[name="method"]').click()
}

const api = isTunnel ? "http://localhost:4001" : ""

const loggedIn = (email) => {
  cy.request(api + "/anything").should((res) => {
    console.log({ body: res.body })

    if (!isTunnel) {
      expect(res.body.headers["Authorization"]).to.not.be.empty
      cy.task(
        "verify",
        res.body.headers["Authorization"].replace(/bearer /gi, ""),
      ).then((decoded) => {
        expect(decoded.session.identity.traits.email).to.equal(email)
      })
    } else {
      expect(res.body.headers["cookie"].indexOf("ory_session_") > -1).to.be.true
    }
  })
}

describe("ory proxy", () => {
  const email = randomEmail()
  const password = randomPassword()
  before(() => {
    cy.clearCookies({ domain: null })
  })

  it("navigation works", () => {
    cy.visit(prefix + "/ui/registration")
    cy.get('[data-testid="cta-link"]').click()
    cy.location("pathname").should("eq", prefix + "/ui/login")
  })

  it("should be able to execute registration", () => {
    cy.visit(prefix + "/ui/registration")
    cy.get('[name="traits.email"]').type(email)
    cy.get('[name="password"]').type(password)
    cy.get('[name="method"]').click()
    if (isTunnel) {
      cy.location("host").should("eq", "localhost:4001")
    }

    loggedIn(email)
  })

  it("should be able to execute login", () => {
    login(email, password)
    if (isTunnel) {
      cy.location("host").should("eq", "localhost:4001")
    }

    loggedIn(email)
  })

  it("should be able to end up locally when signing in with social", () => {
    cy.visit(prefix + "/ui/login")
    cy.get('[value="SnUimsDjTxePInF-"]').click()
    cy.location("host").should("eq", "localhost:4445")

    cy.get('[name="email"]').type("foo@bar.com")
    cy.get('[name="password"]').type("foobar")
    cy.get("#accept").click()

    cy.location("host").should("eq", "localhost:4000")
  })

  it("should be able to execute logout", () => {
    login(email, password)
    loggedIn(email)

    cy.visit(prefix + "/ui/welcome")
    cy.get('[data-testid="logout"] a').should(
      "have.attr",
      "aria-disabled",
      "false",
    )
    cy.get('[data-testid="logout"] a').click()

    if (isTunnel) {
      cy.location("host").should("eq", "localhost:4001")
    }

    cy.request(api + "/anything").should((res) => {
      expect(res.body.headers["Authorization"]).to.be.undefined
    })
  })
})
