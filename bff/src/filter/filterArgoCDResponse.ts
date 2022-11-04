export const filterResponse = (responseBuffer) => {
  const response = responseBuffer.toString("utf8");
  const argocdResponse = JSON.parse(response) || {};
  const { items, metadata } = argocdResponse;
  if (items && Array.isArray(items) && items.length > 0) {
    const mappedItems = items.map((e) => {
      delete e.metadata.labels["argocd.argoproj.io/instance"];
      e.metadata.managedFields = [];
      e.metadata.namespace = "";
      e.metadata.annotations = {};
      (e.status?.operationState?.syncResult?.resources || []).forEach((r) => {
        r["group"] = "";
        r["namespace"] = "";
        r["message"] = "";
      });
      (e.status?.resources || []).forEach((r) => {
        r["group"] = "";
        r["namespace"] = "";
      });
      return e;
    });
    return JSON.stringify({
      items: mappedItems,
      metadata,
    });
  }
  return response;
};
