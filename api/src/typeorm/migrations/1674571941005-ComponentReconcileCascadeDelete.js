module.exports = class ComponentReconcileCascadeDelete1674571941005 {
  name = 'ComponentReconcileCascadeDelete1674571941005';

  async up(queryRunner) {
    await queryRunner.query(
      `ALTER TABLE \`component_reconcile\` DROP FOREIGN KEY \`FK_8d152c6e1ad66defa9d4106585d\``
    );
    await queryRunner.query(
      `ALTER TABLE \`component_reconcile\` ADD CONSTRAINT \`FK_8d152c6e1ad66defa9d4106585d\` FOREIGN KEY (\`componentId\`) REFERENCES \`components\`(\`id\`) ON DELETE CASCADE ON UPDATE NO ACTION`
    );
  }

  async down(queryRunner) {
    await queryRunner.query(
      `ALTER TABLE \`component_reconcile\` DROP FOREIGN KEY \`FK_8d152c6e1ad66defa9d4106585d\``
    );
    await queryRunner.query(
      `ALTER TABLE \`component_reconcile\` ADD CONSTRAINT \`FK_8d152c6e1ad66defa9d4106585d\` FOREIGN KEY (\`componentId\`) REFERENCES \`components\`(\`id\`) ON DELETE NO ACTION ON UPDATE NO ACTION`
    );
  }
};
