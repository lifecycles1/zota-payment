import { useParams } from "react-router-dom";

const Profile = () => {
  const { id } = useParams();

  return (
    <div className="profile-container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">User Profile</h1>
      <p>Welcome to the user profile page.</p>
      <section className="bg-gray-100 py-8">
        <p>UID: {id}</p>
      </section>
    </div>
  );
};

export default Profile;
