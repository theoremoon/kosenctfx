import ApolloClient from "apollo-boost";
import { SERVER_ADDRESS } from "./env";

export const apolloClient = new ApolloClient({
  uri: SERVER_ADDRESS + "query"
});
