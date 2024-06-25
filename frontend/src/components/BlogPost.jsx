import React, { useState, useEffect } from "react";
import axios from "axios";
import { Link } from "react-router-dom";

const BlogPost = () => {
  const [blogData, setBlogData] = useState([]);
  const [loading, setLoading] = useState(false);

  const allBlog = () => {
    setLoading(true);
    axios
      .get(`${process.env.REACT_APP_BACKEND_URL}/api/get-blog`, {
        withCredentials: true,
      })
      .then((response) => {
        setLoading(false);
        setBlogData(response?.data?.data);
        console.log(response?.data?.data);
      })
      .catch((error) => {
        setLoading(false);
        console.log(error);
      });
  };

  useEffect(() => {
    allBlog();
  }, []);

  return (
    <>
      {loading && (
        <div className="text-2xl font-bold text-center px-56 pt-24">
          <h1>LOADING.....</h1>
        </div>
      )}
      <div className="container my-12 mx-auto px-4 md:px-12">
        <div className="flex flex-wrap -mx-1 lg:-mx-4">
          {blogData?.map((blog) => (
            <div
              className="my-1 px-1 w-full md:w-1/2 lg:my-4 lg:px-4 lg:w-1/3"
              key={blog.id}
            >
              <article className="overflow-hidden rounded-lg shadow-lg">
                <Link to={`/detail/${blog.id}`}>
                  <img
                    alt="Blog"
                    className="block h-72 w-full"
                    src={blog?.image}
                  />
                </Link>

                <header className="flex items-center justify-between leading-tight p-2 md:p-4">
                  <h1 className="text-lg">
                    <Link
                      className="no-underline hover:underline text-black"
                      to={`/detail/${blog.id}`}
                    >
                      {blog.title}
                    </Link>
                  </h1>
                  <p className="text-grey-darker text-sm">
                    {new Date(blog?.CreatedAt).toLocaleString()}
                  </p>
                </header>

                <footer className="flex items-center justify-between leading-none p-2 md:p-4">
                  <Link
                    className="flex items-center no-underline hover:underline text-black"
                    to={`/detail/${blog.id}`}
                  >
                    <img
                      alt="User"
                      className="block rounded-full w-5 h-5"
                      src={blog?.user?.profileImage || blog?.image}
                    />
                    <p className="ml-2 text-sm">
                      {blog?.user?.first_name} {blog?.user?.last_name}
                    </p>
                  </Link>
                  <button className="no-underline text-grey-darker hover:text-red-dark">
                    <span className="hidden">Like</span>
                    <i className="fa fa-heart"></i>
                  </button>
                </footer>
              </article>
            </div>
          ))}
        </div>
      </div>
    </>
  );
};

export default BlogPost;
