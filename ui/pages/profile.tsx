import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  VStack,
} from "@chakra-ui/react";
import CountrySelector from "components/countryselector";
import Right from "components/right";
import { GetStaticProps } from "next";
import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import Loading from "../components/loading";
import { api } from "../lib/api";
import useAccount, { Account, fetchAccount } from "../lib/api/account";
import useMessage from "../lib/useMessage";

type UpdateParams = {
  teamname: string;
  password: string;
};

interface ProfileProps {
  account: Account | null;
}

const Profile = ({ account: defaultAccount }: ProfileProps) => {
  const { message, error } = useMessage();
  const { data: account, mutate } = useAccount(defaultAccount);
  const [country, setCountry] = useState(account?.country || "");
  const { register, setValue, handleSubmit } = useForm({
    defaultValues: {
      teamname: account?.teamname,
      password: "",
    },
  });
  const onSubmit: SubmitHandler<UpdateParams> = async (data) => {
    try {
      const res = await api.post("/update-profile", {
        ...data,
        country: country,
      });
      message(res);
      mutate();
      setValue("password", "");
    } catch (e) {
      error(e);
    }
  };

  if (account === undefined) {
    return <Loading />;
  }

  return (
    <Box w="sm" mx="auto" mt="10">
      <form onSubmit={handleSubmit(onSubmit)}>
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

export const getStaticProps: GetStaticProps<ProfileProps> = async () => {
  const account = await fetchAccount();
  return {
    props: {
      account: account,
    },
  };
};

export default Profile;
