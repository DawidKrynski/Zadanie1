describe('Logowanie - formularz', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('powinien wyswietlic formularz logowania', () => {
    cy.get('#user-name').should('be.visible')
    cy.get('#password').should('be.visible')
    cy.get('#login-button').should('be.visible')
    cy.get('.login_logo').should('be.visible')
    cy.get('.login_logo').should('contain', 'Swag Labs')
  })

  it('powinien zalogowac sie poprawnymi danymi', () => {
    cy.get('#user-name').type('standard_user')
    cy.get('#password').type('secret_sauce')
    cy.get('#login-button').click()
    cy.url().should('include', '/inventory.html')
    cy.get('.inventory_list').should('be.visible')
  })

  it('powinien wyswietlic blad przy zlym hasle', () => {
    cy.get('#user-name').type('standard_user')
    cy.get('#password').type('wrong_password')
    cy.get('#login-button').click()
    cy.get('[data-test="error"]').should('be.visible')
    cy.get('[data-test="error"]').should('contain', 'Username and password do not match')
  })
})
