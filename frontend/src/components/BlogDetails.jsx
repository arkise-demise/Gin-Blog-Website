import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import axios from "axios";

const BlogDetail = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [singlePost, setSinglePost] = useState({});

  useEffect(() => {
    const User = localStorage.getItem("user");
    if (!User) {
      navigate("/login");
    }
  }, [navigate]);

  const singleBlog = () => {
    axios
      .get(`${process.env.REACT_APP_BACKEND_URL}/api/get-blog/${id}`, {
        withCredentials: true,
      })
      .then((response) => {
        setSinglePost(response.data.data);
        console.log(response.data.data);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(() => {
    singleBlog();
  }, []);

  return (
    <div className="relative">
      <div className="max-w-3xl mb-10 rounded overflow-hidden flex flex-col mx-auto text-center">
        <div className="max-w-3xl mx-auto text-xl sm:text-4xl font-semibold inline-block hover:text-indigo-600 transition duration-500 ease-in-out inline-block mb-2">
          The Best Activewear from the Nordstrom Anniversary Sale
        </div>

        <img className="w-full h-96 my-4" src={singlePost?.image} alt="Blog" />
        <p className="text-gray-700 text-base leading-8 max-w-2xl mx-auto">
          Author: {singlePost?.user?.first_name} {singlePost?.user?.last_name}
        </p>

        <hr />
      </div>

      <div className="max-w-3xl mx-auto">
        <div className="mt-3 bg-white rounded-b lg:rounded-b-none lg:rounded-r flex flex-col justify-between leading-normal">
          <div className="">
            <p className="text-base leading-8 my-5">{singlePost?.desc}</p>
          </div>
        </div>
      </div>
    </div>
  );
};
export default BlogDetail;
