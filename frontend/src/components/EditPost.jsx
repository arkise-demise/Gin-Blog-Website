import React, { useState, useEffect } from "react";
import axios from "axios";
import { useForm } from "react-hook-form";
import { useParams, useNavigate } from "react-router-dom";
import { useSnackbar } from "react-simple-snackbar";

const EditPost = () => {
  const [singlePost, setSinglePost] = useState();
  const [loading, setLoading] = useState(false);
  const { id } = useParams();
  const navigate = useNavigate();
  const [image, setImage] = useState();
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm();

  const options = {
    position: "bottom-right",
    style: {
      backgroundColor: "gray",
      border: "2px solid lightgreen",
      color: "white",
      fontFamily: "Menlo, monospace",
      fontSize: "20px",
      textAlign: "center",
    },
    closeStyle: {
      color: "lightcoral",
      fontSize: "16px",
    },
  };

  const [openSnackbar] = useSnackbar(options);

  const fetchSingleBlog = () => {
    axios
      .get(`${process.env.REACT_APP_BACKEND_URL}/api/get-blog/${id}`, {
        withCredentials: true,
      })
      .then((response) => {
        setSinglePost(response.data.data);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(() => {
    const User = localStorage.getItem("user");
    if (!User) {
      navigate("/login");
    } else {
      fetchSingleBlog();
    }
  }, [navigate, id]);

  const onSubmit = (data) => {
    setLoading(true);
    const body = {
      ...data,
      image: singlePost?.image,
    };

    axios
      .put(
        `${process.env.REACT_APP_BACKEND_URL}/api/update-blog/${id}`,
        { ...body },
        {
          withCredentials: true,
        }
      )
      .then((response) => {
        openSnackbar("Post Updated Successfully");
        setLoading(false);
        navigate("/personal");
      })
      .catch((error) => {
        openSnackbar("Oops!, Post is not updated");
        setLoading(false);
      });
  };

  return (
    <div className="max-w-screen-md mx-auto p-5">
      <form className="w-full" onSubmit={handleSubmit(onSubmit)}>
        <div className="flex flex-wrap -mx-3 mb-6">
          <div className="w-full md:w-full px-3 mb-6 md:mb-0">
            <label
              className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2"
              htmlFor="grid-first-name"
            >
              Title
            </label>
            <input
              className="appearance-none block w-full bg-gray-200 text-gray-700 border border-red-500 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white"
              id="grid-first-name"
              type="text"
              placeholder="title"
              name="title"
              autoComplete="off"
              defaultValue={singlePost?.title}
              {...register("title", { required: true })}
            />
            {errors.title && errors.title.type === "required" && (
              <p className="text-red-500 text-xs italic">
                Please fill out this field.
              </p>
            )}
          </div>
        </div>
        <div className="flex flex-wrap -mx-3 items-center lg:items-start mb-6">
          <div className="w-full px-3">
            <label title="click to select a picture">
              <input
                accept="image/*"
                className="hidden"
                id="banner"
                type="file"
                name="image"
                // onChange={handleImage}
                visibility="hidden"
              />
              <div className="flex flex-col">
                <div className="pb-2">Upload Image</div>
                {image || singlePost ? (
                  <div className="pt-4">
                    <img
                      className="-object-contain -mt-8 p-5 w-1/2"
                      src={image ? image.image : singlePost?.image}
                      alt=""
                    />
                  </div>
                ) : (
                  <div className="pb-5">
                    <img
                      src="/upload-image.svg"
                      style={{ background: "#EFEFEF" }}
                      className="h-full w-48"
                    />
                  </div>
                )}
              </div>
            </label>
          </div>
          <div className="flex items-center justify-center px-5">
            <button
              className="shadow bg-indigo-600 hover:bg-indigo-400 focus:shadow-outline focus:outline-none text-white font-bold py-2 px-6 rounded"
              // onClick={uploadImage}
            >
              Upload Image
            </button>
          </div>
        </div>
        <div className="flex flex-wrap -mx-3 mb-6">
          <div className="w-full px-3">
            <label
              className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2"
              htmlFor="grid-password"
            >
              Description
            </label>
            <textarea
              rows="10"
              name="desc"
              defaultValue={singlePost?.desc}
              className="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
              {...register("desc", { required: true })}
            ></textarea>
            {errors.desc && errors.desc.type === "required" && (
              <p className="text-red-500 text-xs italic">
                Please fill out this field.
              </p>
            )}
          </div>
          <div className="flex justify-between w-full px-3">
            <button
              className="shadow bg-indigo-600 hover:bg-indigo-400 focus:shadow-outline focus:outline-none text-white font-bold py-2 px-6 rounded"
              type="submit"
              disabled={loading ? true : false}
            >
              {loading ? "Loading" : "Update Post"}
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default EditPost;
