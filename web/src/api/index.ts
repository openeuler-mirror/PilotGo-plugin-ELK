import request from "./request";

// test
export function testHttps() {
  return request({
    url: "/plugin/elk/api/search",
    method: "post",
  });
}