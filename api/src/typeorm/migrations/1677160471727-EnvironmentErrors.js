module.exports = class EnvironmentErrors1677160471727 {
    name = "EnvironmentErrors1677160471727";

    async up(queryRunner) {
        await queryRunner.query(
            `ALTER TABLE \`environment\` ADD \`errorType\` INT DEFAULT NULL`
          );
        await queryRunner.query(
            `ALTER TABLE \`environment\` ADD \`errorMessage\` json DEFAULT NULL`
          );
    }

    async down(queryRunner) {
        await queryRunner.query(
            `ALTER TABLE \`environment\` DROP COLUMN \`errorType\``
          );
          await queryRunner.query(
            `ALTER TABLE \`environment\` DROP COLUMN \`errorMessage\``
          );
    }

}
