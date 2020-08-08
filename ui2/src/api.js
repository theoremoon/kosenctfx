import axios from "axios";
import { SERVER_ADDRESS } from "./env";

const API = axios.create({
  baseURL: SERVER_ADDRESS,
  withCredentials: true
});
export default API;
