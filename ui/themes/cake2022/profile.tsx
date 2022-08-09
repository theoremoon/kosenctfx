import { ProfileProps } from "props/profile";
import CountrySelector from "./components/countrySelector";
import Input from "./components/input";
import Button from "./components/button";
import Right from "./components/right";

const Profile = ({ register, onSubmit, country, setCountry }: ProfileProps) => {
  return (
    <form onSubmit={onSubmit} className="input-form">
      <div className="form-item">
        <label htmlFor="teamname">teamname</label>
        <Input
          id="teamname"
          className="form-input"
          autoComplete="username"
          {...register("teamname", { required: true })}
        ></Input>
      </div>
      <div className="form-item">
        <label htmlFor="password">password</label>
        <Input
          id="password"
          type="password"
          className="form-input"
          autoComplete="new-password"
          {...register("password", { required: false })}
        ></Input>
      </div>
      <div className="form-item">
        <label htmlFor="password">password</label>
        <label htmlFor="country">country</label>
        <CountrySelector
          id="country"
          className="form-input"
          value={country}
          onChange={(e) => setCountry(e.target.value)}
        />
      </div>

      <div className="form-item">
        <Right>
          <Button type="submit">Update</Button>
        </Right>
      </div>
    </form>
  );
};

export default Profile;
