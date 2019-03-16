function (user, context, callback) {
    var whitelist = ['example1@gmail.com', 'example2@gmail.com']; //authorized emails
    var userHasAccess = whitelist.some(
      function (email) {
        return user.email.toLowerCase() === email;
      });
    const is_social = context.connectionStrategy === context.connection;
    if (!userHasAccess && is_social) {
      return callback(new UnauthorizedError('Access denied.'));
    }

    return callback(null, user, context);
}
