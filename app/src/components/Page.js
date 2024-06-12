import React from 'react';
import { Outlet } from 'react-router-dom';
import MenuBar from './MenuBar';

const Page = () => {
  return (
    <div>
      <MenuBar />
      <Outlet />
    </div>
  );
}

export default Page;
