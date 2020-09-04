import md5 from "blueimp-md5";

export default function colorhash(x) {
  return "#" + md5(x).substr(-6);
}
