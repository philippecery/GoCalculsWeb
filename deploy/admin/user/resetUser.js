try {
    conn = new Mongo(uri);
    db = conn.getDB("maths");
    db.auth("maths_rw", passwordPrompt());
    db.users.deleteMany(
        { userid: "admin" }
    );
    db.users.insertOne(
        {
            userid: "admin",
            emailaddress: "admin@maths.com",
            password: adminPassword,
            firstname: "Admin",
            lastname: "Admin",
            role: 1,
            status: 2
        }
    );
    if (db.users.find().count() == 1) {
        print("Admin user is created/reset");
    }
} catch (err) {
    print("Admin user could not be created. Cause: " + err);
}
