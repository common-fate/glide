const regexSuffixless = /\/[^/.]+$/; // e.g. "/some/page" but not "/", "/some/" or "/some.jpg"
const regexTrailingSlash = /.+\/$/; // e.g. "/some/" or "/some/page/" but not root "/"
const dynamicRouteRegex = /\/subpath\/\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b/; // e.g /urs/some-uuid; // e.g. '/subpath/uuid'

const commonFateIdRegex = /\w{3}_\w{27}/;

export const fixURI = (uri: string): string => {
  uri = uri.replace(commonFateIdRegex, "[id]");

  //Checks for dynamic route and retrieves the proper [id].html file
  if (uri.match(dynamicRouteRegex)) {
    return "/subpath/[id].html";
  }

  // Append ".html" to origin request
  if (uri.match(regexSuffixless)) {
    return uri + ".html";
  }

  // Append "index.html" to origin request
  if (uri.match(regexTrailingSlash)) {
    return uri + "index.html";
  }

  return uri;
};
