function getStockData(id) {
  element = document.getElementById(`stock-info-panel-${id}`);
  if (!element) {
    return "";
  }
  return JSON.parse(
    document.getElementById(`stock-info-panel-${id}`).dataset.stock,
  );
}
