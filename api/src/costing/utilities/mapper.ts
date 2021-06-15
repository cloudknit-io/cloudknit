export class Mapper {
  static Map<T1, T2>(mapFrom: T1, mapTo: T2): T2 {
    for (let key of Object.keys(mapFrom)) {
      if (key in mapTo) {
        mapTo[key] = mapFrom[key]
      }
    }
    return mapTo
  }
}
