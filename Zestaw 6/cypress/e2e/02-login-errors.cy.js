describe('Logowanie - bledy', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('powinien wyswietlic blad przy pustych polach', () => {
    cy.get('#login-button').click()
    cy.get('[data-test="error"]').should('be.visible')
    cy.get('[data-test="error"]').should('contain', 'Username is required')
  })

  it('powinien wyswietlic blad dla zablokowanego uzytkownika', () => {
    cy.get('#user-name').type('locked_out_user')
    cy.get('#password').type('secret_sauce')
    cy.get('#login-button').click()
    cy.get('[data-test="error"]').should('be.visible')
    cy.get('[data-test="error"]').should('contain', 'locked out')
  })
})
