import React from "react";

type InputProps = Omit<React.ComponentPropsWithRef<"input">, "className">;

const Input = React.forwardRef((props: InputProps, ref) => {
  return (
    <input
      {...props}
      ref={ref}
      className="w-full bg-transparent color-black border-b focus:border-b-2 mb-px focus:mb-0 border-pink-600 outline-none text-center text-xl"
    />
  );
});

export default Input;
