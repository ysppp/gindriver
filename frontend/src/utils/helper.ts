export function queryURLParams(url: string) {
  let obj: {[props: string]: string} = {};
  url.replace(/([^?=&#]+)=([^?=&#]+)/g, (...[, $1, $2]) => obj[<string>$1] = <string>$2);
  return obj;
}