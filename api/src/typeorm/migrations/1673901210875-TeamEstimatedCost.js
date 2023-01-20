module.exports = class TeamEstimatedCost1673901210875 {
  name = 'TeamEstimatedCost1673901210875';

  async up(queryRunner) {
    await queryRunner.query(
      `ALTER TABLE \`team\` ADD \`estimated_cost\` decimal(10,3) NOT NULL DEFAULT '0.000'`
    );
  }

  async down(queryRunner) {
    await queryRunner.query(
      `ALTER TABLE \`team\` DROP COLUMN \`estimated_cost\``
    );
  }
};
