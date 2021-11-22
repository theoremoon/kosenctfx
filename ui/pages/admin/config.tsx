import { css, jsx } from "@emotion/react";
import { pink } from "../../lib/color";

const Config = () => {
  return (
    <div className="center-container">
      <div className="w-full lg:w-1/2 mx-auto border-2 border-pink-600 rounded p-10">
        <form>
          <h2 className="text-xl font-bold">Bucket Settings</h2>
          <div>
            <label className="font-bold block">Bucket Endpoint</label>
            <input type="text" className="w-full bg-transparent border-b-1" />
          </div>
        </form>
      </div>
    </div>
  );
};

export default Config;
