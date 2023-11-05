import { LoginProps } from "props/login";
import React from "react";
import Input from "./components/input";
import Button from "./components/button";
import Link from "next/link";
import Right from "./components/right";

const Login = ({ register, onSubmit }: LoginProps) => {
  return (
    <form onSubmit={onSubmit} className="input-form">
      <div className="form-item">
        <label style={{ display: "block" }} htmlFor="teamname">
          teamname
        </label>
        <Input
          className="form-input"
          id="teamname"
          placeholder="yoshiking"
          {...register("teamname", { required: true })}
        />
      </div>
      <div className="form-item">
        <label htmlFor="password">password</label>
        <Input
          className="form-input"
          id="password"
          type="password"
          {...register("password", { required: true })}
        />
      </div>
      <div className="form-item">
        <Right>
          <Button type="submit">Login</Button>
        </Right>
      </div>
      <div>
        <Right>
          <Link href="/passwordreset_request">
            Forgot your password? You can reset your password here.
          </Link>
        </Right>
      </div>
    </form>
  );
};

export default Login;
