package gologix

import "fmt"

func (client *Client) ArbitraryMessage(service CIPService, path Serializable, SendData Serializable) ([]cipItem, error) {

	reqitems := make([]cipItem, 2)
	//reqitems[0] = cipItem{Header: cipItemHeader{ID: cipItem_Null}}
	reqitems[0] = NewItem(cipItem_ConnectionAddress, &client.OTNetworkConnectionID)

	readmsg := msgCIPConnectedServiceReq{
		SequenceCount: client.Sequencer(),
		Service:       service,
		PathLength:    byte(path.Len() / 2),
	}

	// setup item
	reqitems[1] = NewItem(cipItem_ConnectedData, readmsg)
	// add path
	reqitems[1].Marshal(path.Bytes())
	// add service specific data
	//number_of_attr_to_receive := 3
	//attr1_symbol_name := 1
	//attr2_symbol_type := 2
	//attr8_arraydims := 8
	//reqitems[1].Marshal([4]uint16{uint16(number_of_attr_to_receive), uint16(attr1_symbol_name), uint16(attr2_symbol_type), uint16(attr8_arraydims)})
	reqitems[1].Marshal(SendData.Bytes())
	reqitems[1].Marshal([3]uint16{1, 0, 1})

	hdr, data, err := client.send_recv_data(cipCommandSendUnregistered, MarshalItems(reqitems))
	if err != nil {
		return nil, err
	}
	_ = hdr
	_ = data
	//data_hdr := ListInstanceHeader{}
	//binary.Read(data, binary.LittleEndian, &data_hdr)

	// first six bytes are zero.
	padding := make([]byte, 6)
	_, err = data.Read(padding)
	if err != nil {
		return nil, fmt.Errorf("douldn't read data. %w", err)
	}

	resp_items, err := ReadItems(data)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse items. %w", err)
	}

	return resp_items, nil
}
