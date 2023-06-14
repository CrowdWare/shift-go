# Todo

- Database for message storage (StorJ, libp2p dht)
    On Storj you have 25 GB for free to store your own data. The user could create a bucket for shift and donate it to the shift movement.
    Storj is decentral and more secure than libp2p.
    
- No service has to run in background, because we are not going P2P.
- Scooping is done all 20 hours, user just has to start it on a daily base.
- Giving is now as follows:
    Receiver shows QR code or sends a hex code with a proposal (public key, amount, purpose, isOnline=can ping storj)
    Giver accepts the proposal and stores a record for the receiver on storj if isOnline is set.
    Receiver pulls the record from storj.

    If giver or receiver is offline we will fallback to QR code transmission.  
    Alternatively instead of storing a record on storj the giver can create a QR code and receiver scan that record.
    In both cases sending a record or creating a QR code withdraws the amount from giver account and scanning or downloading the record books the amount to receivers account. The record is only useful for the receiver and will loose the validity after x days.

    In the case that a record cannot be pulled by the receiver. The receiver has to inform the giver. Than the giver pulls the record back from storj and the amount is booked to the account.

    A record has the following key:   public key of the receiver - "LMP" - public key of giver

    QR code first !!!

- REST calls are optional, so when REST call fails, the client will continue to operate, no warnings will be shown
    Invite code, name, etc are optional the user does not have to enter these values
    