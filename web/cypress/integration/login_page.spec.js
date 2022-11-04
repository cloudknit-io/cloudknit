describe('Login page', () => {
	before(() => {
		cy.visit('/');
	});

	describe('as an unauthenticated user', () => {
		it('successfully loads', function () {
			cy.get('.login-form').should('exist');
		});
	});
});
