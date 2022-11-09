import { ComponentAudit } from "src/reconciliation/dtos/componentAudit.dto";
import { EnvironmentAudit } from "src/reconciliation/dtos/environmentAudit.dto";
import { ComponentReconcileDto } from "src/reconciliation/dtos/reconcile.Dto";
import { Organization } from "src/typeorm";
import { Component } from "src/typeorm/costing/Component";
import { ComponentReconcile } from "src/typeorm/reconciliation/component-reconcile.entity";
import { EnvironmentReconcile } from "src/typeorm/reconciliation/environment-reconcile.entity";
import { CostComponent, Resource } from "src/typeorm/resources/Resource.entity";

export class Mapper {
  static Map<T1, T2>(mapFrom: T1, mapTo: T2): T2 {
    for (let key of Object.keys(mapFrom)) {
      if (key in mapTo) {
        mapTo[key] = mapFrom[key];
      }
    }
    return mapTo;
  }

  static getStreamData(mapFrom: Component[]): {} {
    const data = {};
    const teams = [...new Set(mapFrom.map((e) => e.teamName))];
    const environments = [
      ...new Set(mapFrom.map((e) => `${e.teamName}$$$${e.environment.name}`)),
    ];

    data["teams"] = teams.map((t) => ({
      teamId: t,
      teamName: t,
      cost: mapFrom
        .filter((e) => e.teamName === t)
        .reduce((p, c, _i) => p + Number(c.cost), 0),
    }));

    data["environments"] = environments.map((e) => {
      const [teamName, environmentName] = e.split("$$$");
      return {
        environmentId: e.replace("$$$", "-"),
        environmentName,
        cost: mapFrom
          .filter(
            (e) =>
              e.teamName === teamName && e.environment.name === environmentName
          )
          .reduce((p, c, _i) => p + Number(c.cost), 0),
      };
    });

    data["components"] = mapFrom.map((e) => ({
      componentId: e.id,
      componentName: e.componentName,
      cost: Number(e.cost),
    }));

    return data;
  }

  static getResource(data: any) {
    return {
      name: data["name"],
      hourlyCost: data["hourlyCost"],
      monthlyCost: data["monthlyCost"],
      subresources: [],
      resourceName: data.resourceName,
      costComponents: [],
    };
  }

  static getCostComponent(data: any) {
    return {
      hourlyCost: data["hourlyCost"],
      hourlyQuantity: data["hourlyQuantity"],
      monthlyCost: data["monthlyCost"],
      monthlyQuantity: data["monthlyQuantity"],
      name: data["name"],
      price: data["price"],
      unit: data["unit"],
    };
  }

  static mapToResourceEntity(
    org: Organization,
    component: Component,
    resources: any[] = [],
    parentId: string = null
  ): Resource[] {
    if (resources.length === 0) return [];
    return resources.map((resource) => {
      const id = `${component.id}-${resource.name}`;
      return {
        organization: org,
        id: id,
        name: resource.name,
        componentId: component.id,
        hourlyCost: resource.hourlyCost,
        monthlyCost: resource.monthlyCost,
        parentId: parentId,
        component: component,
        subresources: this.mapToResourceEntity(
          org,
          component,
          resource.subresources,
          id
        ),
        costComponents: this.mapToCostComponentEntity(
          org,
          id,
          resource.costComponents
        ),
      };
    });
  }

  static mapToCostComponentEntity(
    org: Organization,
    id: string,
    costComponents: any[] = []
  ): CostComponent[] {
    if (costComponents.length === 0) return [];
    return costComponents.map((cc) => ({
      organization: org,
      id: `${id}-${cc.name}`,
      name: cc.name,
      hourlyCost: cc.hourlyCost,
      hourlyQuantity: cc.hourlyQuantity,
      monthlyCost: cc.monthlyCost,
      monthlyQuantity: cc.monthlyQuantity,
      price: cc.price,
      unit: cc.unit,
    }));
  }

  static mapToComponentReconcile(
    org: Organization,
    environmentReconcile: EnvironmentReconcile,
    componentReconciles: ComponentReconcileDto[]
  ): ComponentReconcile[] {
    const mappedData: ComponentReconcile[] = componentReconciles.map((cr) => ({
      reconcile_id: parseInt(cr.reconcileId) || null,
      name: cr.name,
      environmentReconcile: environmentReconcile,
      start_date_time: cr.startDateTime,
      status: cr.status,
      end_date_time: cr.endDateTime,
      organization: org
    }));

    return mappedData;
  }

  static getComponentAuditList(
    components: ComponentReconcile[]
  ): ComponentAudit[] {
    return components.map((c) => {
      let diff = -1;
      if (c.end_date_time) {
        const ed = new Date(c.end_date_time).getTime();
        const sd = new Date(c.start_date_time).getTime()
        diff = ed - sd;
      }

      return {
        reconcileId: c.reconcile_id,
        startDateTime: c.start_date_time,
        duration: diff,
        status: c.status,
        approvedBy: c.approved_by,
      };
    });
  }

  static getEnvironmentAuditList(
    components: EnvironmentReconcile[]
  ): EnvironmentAudit[] {
    return components.map((c) => {
      let diff = -1;
      if (c.end_date_time) {
        const ed = new Date(c.end_date_time).getTime();
        const sd = new Date(c.start_date_time).getTime()
        diff = ed - sd;
      }

      return {
        reconcileId: c.reconcile_id,
        startDateTime: c.start_date_time,
        duration: diff,
        status: c.status,
      };
    });
  }
}