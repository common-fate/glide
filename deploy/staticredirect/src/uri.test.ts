import { fixURI } from "./uri";

describe("fixURI", () => {
  it("should add index.html to a trailing slash", () => {
    expect(fixURI("/admin/")).toEqual("/admin/index.html");
  });
  it("should add .html to a page", () => {
    expect(fixURI("/admin/access-rules")).toEqual("/admin/access-rules.html");
  });
  it("should add [id] to a Common Fate ID", () => {
    expect(
      fixURI("/admin/access-rules/rul_29kaLgLmxb7b8rAcy4YuE9bROTx")
    ).toEqual("/admin/access-rules/[id].html");
  });
  it("should return the URI if nothing needs changing", () => {
    expect(fixURI("/index.html")).toEqual("/index.html");
  });
});
