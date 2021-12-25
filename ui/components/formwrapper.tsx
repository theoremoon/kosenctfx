interface FormWrapperProps {
  children: React.ReactNode;
}

const FormWrapper = ({ children }: ReactWrapperProps) => {
  return (
    <div className="w-full lg:w-3/4 xl:w-1/2 mx-auto border-2 border-pink-600 rounded p-10 mt-20">
      {children}
    </div>
  );
};

export default FormWrapper;
