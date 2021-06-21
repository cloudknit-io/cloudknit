import { Component } from 'src/typeorm/costing/entities/Component'
import { Resource } from 'src/typeorm/resources/Resource.entity'

export class Mapper {
  static Map<T1, T2>(mapFrom: T1, mapTo: T2): T2 {
    for (let key of Object.keys(mapFrom)) {
      if (key in mapTo) {
        mapTo[key] = mapFrom[key]
      }
    }
    return mapTo
  }

  static getStreamData(mapFrom: Component[]): {} {
    const data = {}
    const teams = [...new Set(mapFrom.map((e) => e.teamName))]
    const environments = [
      ...new Set(mapFrom.map((e) => `${e.teamName}$$$${e.environmentName}`)),
    ]
    data['teams'] = teams.map((t) => ({
      teamId: t,
      teamName: t,
      cost: mapFrom
        .filter((e) => e.teamName === t)
        .reduce((p, c, _i) => p + Number(c.cost), 0),
    }))

    data['environments'] = environments.map((e) => {
      const [teamName, environmentName] = e.split('$$$')
      return {
        environmentId: e.replace('$$$', '-'),
        environmentName,
        cost: mapFrom
          .filter(
            (e) =>
              e.teamName === teamName && e.environmentName === environmentName,
          )
          .reduce((p, c, _i) => p + Number(c.cost), 0),
      }
    })

    data['components'] = mapFrom.map((e) => ({
      componentId: e.id,
      componentName: e.componentName,
      cost: Number(e.cost),
    }))

    return data
  }

  static getResource(data: any) {
    return {
      name: data['name'],
      hourlyCost: data['hourlyCost'],
      monthlyCost: data['monthlyCost'],
      subresources: [],
      resourceName: data.resourceName,
    };
  }
}
