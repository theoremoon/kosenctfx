import {
  Box,
  VStack,
  FormControl,
  FormLabel,
  Input,
  Button,
} from "@chakra-ui/react";
import { ProfileProps } from "props/profile";
import Right from "./components/right";
import CountrySelector from "./components/countrySelector";

const Profile = ({ register, onSubmit, country, setCountry }: ProfileProps) => {
  return (
    <Box w="sm" mx="auto" mt="10">
      <form onSubmit={onSubmit}>
        <VStack>
          <FormControl>
            <FormLabel htmlFor="teamname">teamname</FormLabel>
            <Input
              id="teamname"
              variant="flushed"
              autoComplete="username"
              {...register("teamname", { required: true })}
            ></Input>
          </FormControl>
          <FormControl>
            <FormLabel htmlFor="password">password</FormLabel>
            <Input
              id="password"
              type="password"
              autoComplete="new-password"
              {...register("password", { required: false })}
            ></Input>
          </FormControl>
          <FormControl>
            <FormLabel htmlFor="country">country</FormLabel>
            <CountrySelector
              id="country"
              value={country}
              onChange={(e) => setCountry(e.target.value)}
            />
          </FormControl>

          <FormControl>
            <Right>
              <Button type="submit">Update</Button>
            </Right>
          </FormControl>
        </VStack>
      </form>
    </Box>
  );
};

export default Profile;
