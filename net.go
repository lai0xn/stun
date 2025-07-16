package stun
//
// func sendBindingResponse(p *Packet) error {
//   msgHeader := Header{
//     Type: BindingResponse,
//     Length: XORMappedAddressLength,
//     MagicCookie: p.message.Header.MagicCookie,
//     TransactionID: p.message.Header.TransactionID,
//   }
//
//   var msgAttrs Attributes
//   xorAddr,err := serializeAddr(XorMappedAddr{
//     Family: IPV4,
//     IP: p.remoteIP,
//     Port: p.remotePort,
//   }) 
//
//   if err != nil {
//     return err
//   }
//
//   msgAttrs = append(msgAttrs, Attribute{
//     Length: XORMappedAddressLength,
//     Type: XORMappedAddress,
//     Value: xorAddr,
//   })  
//
//   msg := Message{
//     Header: msgHeader,
//     Attributes: msgAttrs,
//   }
//   msgBytes:= msg.Encode()
//
//   p.Write(msgBytes)
//   return nil
// }
