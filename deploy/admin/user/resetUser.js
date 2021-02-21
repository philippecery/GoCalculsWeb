try {
    conn = new Mongo("localhost:27017");
    db = conn.getDB("maths");
    db.auth("maths_rw", passwordPrompt());
    db.users.deleteMany(
        { userid: adminId }
    );
    db.users.insertOne(
        {
            userid: adminId,
            emailaddress: adminEmail,
            password: adminPassword,
            role: 4,
            status: 2
        }
    );
    if (db.users.find().count() == 1) {
        print("Admin user is created/reset");
    }
} catch (err) {
    print("Admin user could not be created. Cause: " + err);
}
