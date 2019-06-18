export function clearSearchContent(searchContent) {
    for (let name in searchContent) {
        if (name !== "pageNumber" && 
        name !== "pageSize" && 
        name != "pageNum" && 
        name != "timeOrder" && 
        name != "reportOrder") {
            searchContent[name] = "";
        }
    }
    return searchContent;
}